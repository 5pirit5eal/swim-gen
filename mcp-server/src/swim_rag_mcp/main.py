import asyncio
import logging
import os

import httpx
from dotenv import load_dotenv
from fastmcp import FastMCP
from fastmcp.exceptions import ToolError
from fastmcp.server.dependencies import get_context
from fastmcp.server.middleware.error_handling import RetryMiddleware
from fastmcp.server.middleware.logging import StructuredLoggingMiddleware
from fastmcp.server.middleware.rate_limiting import (
    RateLimitingMiddleware,
    SlidingWindowRateLimitingMiddleware,
)

from swim_rag_mcp.auth import IdTokenManager
from swim_rag_mcp.schemas import ExportResponse, QueryRequest, QueryResponse

load_dotenv(".config.env")

URL = os.getenv("SWIM_RAG_API_URL", "http://localhost:8080")
MCP_URL = os.getenv("SWIM_RAG_MCP_URL", "http://localhost:5000")
print(f"Using Swim RAG API URL: {URL}")

# Instantiate the manager once at application startup
token_manager = IdTokenManager()


async def make_authenticated_request(
    url: str,
    method: str = "GET",
    json_data: dict = None,  # type: ignore[assignment]
) -> httpx.Response:
    """Make an authenticated request to the Swim RAG API with retries on 401 Unauthorized."""
    ctx = get_context()

    if not json_data:
        json_data = {}

    max_retries = 2
    for attempt in range(max_retries):
        try:
            id_token = await token_manager.get_id_token(URL)  # Use the manager
            headers = {
                "Authorization": f"Bearer {id_token}",
                "Content-Type": "application/json",
            }
            await ctx.debug(
                f"Attempt {attempt + 1}: Making {method} request to {url} with httpx..."
            )
            async with httpx.AsyncClient() as client:
                response = await client.request(
                    method, url, headers=headers, json=json_data, timeout=60
                )
                response.raise_for_status()
            return response
        except httpx.HTTPStatusError as e:
            if e.response.status_code == 401 and attempt < max_retries - 1:
                await ctx.info(
                    "Request failed with 401 Unauthorized. Invalidating cached token and retrying..."
                )
                await token_manager.invalidate_token()  # Invalidate via manager
            else:
                await ctx.error(f"Request failed: {e}")
                raise
        except Exception as e:
            await ctx.error(f"An unexpected error occurred: {e}")
            raise

    raise RuntimeError("Max retries exceeded for authenticated request.")


mcp: FastMCP = FastMCP(
    name="swim-rag-mcp",
    instructions="""
        This is the MCP Server connected to the Swim RAG backend, an application meant for generating and
        exporting german training plans for swimming. It allows the user to query for a personalized training plan,
        edit it and send the edited plan to the Swim RAG backend for export to a PDF file.
    """,
    exclude_tags={"internal"},
    include_tags={"public"},
    on_duplicate_tools="error",  # Handle duplicate registrations
    on_duplicate_resources="warn",
    on_duplicate_prompts="replace",
    middleware=[
        StructuredLoggingMiddleware(
            include_payloads=True, log_level=logging.INFO
        ),
        RateLimitingMiddleware(burst_capacity=20),
        SlidingWindowRateLimitingMiddleware(
            max_requests=100,
            window_minutes=1,
        ),
        RetryMiddleware(max_retries=3),
    ],
)


@mcp.tool(tags={"public"})
async def generate_or_choose_plan(query: QueryRequest) -> QueryResponse:
    """Query the Swim RAG system with a given german query string.
    It parses the request, queries the RAG, generating or choosing a plan, and returns the result as JSON.
    This function can be used to generate a new training plan or choose an existing one from the database.
    It is not suited for editing plans, but rather for querying the Swim RAG backend.
    """
    # Send the request to the Swim RAG backend
    try:
        response = await make_authenticated_request(
            url=URL + "/query",
            method="POST",
            json_data=query.model_dump(),
        )
        response.raise_for_status()  # Raise an error for bad responses
    except httpx.RequestError as e:
        raise ToolError(
            "The connection to the MCP Server is successful, "
            f"but the request to the connected API failed with error: {e}"
        )
    except httpx.HTTPStatusError as e:
        raise ToolError(
            f"HTTP error: {e.response.status_code} - {e.response.text}"
        )
    except Exception as e:
        raise ToolError(f"An unexpected error occurred: {e}")
    # Return the response
    return QueryResponse.model_validate_json(response.text)


@mcp.tool(tags={"public"})
async def export_plan(plan: QueryResponse) -> ExportResponse:
    """Export a plan as a PDF file for easier printing and sharing."""
    try:
        response = await make_authenticated_request(
            url=URL + "/export-pdf",
            method="POST",
            json_data=plan.model_dump(),
        )
        response.raise_for_status()  # Raise an error for bad responses
    except httpx.RequestError as e:
        raise ToolError(f"Request error: {e}")
    except httpx.HTTPStatusError as e:
        raise ToolError(
            f"HTTP error: {e.response.status_code} - {e.response.text}"
        )
    except Exception as e:
        raise ToolError(f"An unexpected error occurred: {e}")

    # Return the response
    return ExportResponse.model_validate_json(response.text)


@mcp.tool(tags={"internal"})
async def scrape_plans_from_web(url: str) -> str:
    """Scrape plans from a given URL.

    Args:
        url (str): The URL to scrape plans from.

    Returns:
        str: Confirmation message of the scraping operation.
    """
    # Here you would implement the logic to scrape plans from
    return f"Plans scraped from {url} successfully."


if __name__ == "__main__":
    asyncio.run(mcp.run_async(transport="http"))
