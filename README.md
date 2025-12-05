# Swim Training Plan Generator

A web application for generating, recommeding and sharing training plans for swimmers, triathletes and trainers.

## Tech Stack

The web application is hosted in Google Cloud. The following components are involved:

1. Identity-Aware Proxy (IAP) to authenticate and identify users
2. Single-Page Frontend hosted in Cloud Run implemented in **Vue.js**
3. Backend-for-Frontend (BFF) hosted in Cloud Run implemented in **Node.js**
4. Backend hosted in Cloud Run implemented in **Go** and **Langchaingo**
5. PostgreSQL database (managed) with **pgvector**

```plaintext
+--------------+       +-------+        +-----------+         +----------+
|              | HTTPS |       | HTTPS  |           |         |          |
|   FRONTEND   +------>|  BFF  +------->|  BACKEND  +-------->| SUPABASE |
|   (Vue.js)   |       |       |        |   (Go)    |         |          |
+------+-------+       +-------+        +-----------+         +----------+
       |                                                           ^     
       |                                                           |     
       |                                                           |     
       |                                                           |     
       +-----------------------------------------------------------+     
```

## Roadmap

### V1: The Prototype

The initial version only contains the core concept of being able to generate individual training plans based on user input (text/form).

- **Generate TP based on user input (free form and configuration menu)**
  - Frontend displays simple one page app with text box and advanced settings
  - Backend provides endpoint for retrieval-augmented generation
  - Database contains training plans scraped from the web (tiny scope) and their embeddings
- **TP are exportable/printable as PDF**
  - Frontend includes Button for exporting PDF
  - Backend allows the user to get PDF exports of Trainingplan data
- **Requests are anonym and no user identification is necessary**
  - Authentication is omitted for v1
  - Any user information is stored encrypted if not configured otherwise
  - Frontend displays the option to not upload results in advanced settings
- **Backend allows user to upload training plans**
  - Backend-only implementation of training plan donation
  - Backend-only implementation of feedback

### V2: User-specific TPs and TP History and Multimodal Generation

The next version takes lessons learned in v1 into consideration and adds user authentication and generation history, making the web application a lot more useful.

- **Add chat interaction for logged-in users for adaptation of results**
  - Builds upon the content of v1 by adding chat-like interaction
  - Chat interaction is integrated as part of the go backend
  - Chat history and user data saved in DB
- **TP can be recommended/donated**
  - Form added to frontend to input/upload your favorite training plans
  - Allow user feedback for plans
  - Endpoint in backend for saving feedback
  - Notation/Abbreviations of different sources can be processed correctly
- **User authentication and authorization**
  - Frontend extended with `Login` page and user info about benefits of logging in
  - Google Auth Platform for user authentication
  - User Email and Password as another option
  - Users get up to 1000 interactions to avoid abuse
  - Whitelisted users with unlimited interactions
- **TP History can be viewed**
  - TPs are saved associated with user
  - Frontend is extended with `History` page to view previously created TPs
- **TP can be shared via URL**
- **Enable deletion of user data**
  - Add url to delete all user data to the user profile
- **Add content required by german/eu law**
  - Impressum
  - Datenschutzerklärung
- **Login possible via Supabase**
- **Add anonymous/guest mode**
  - An anonymous mode is added which enables new users to use the application in the same way as v1.
- **Upload TP based on multimodal user input (pdf, image of hand-written notes)**

### V3: Community-Expansion and Monetization

- **Tutorial page for understanding plans**
  - References input into the plan display with explainations and tutorials
- **Add a public leaderboard and public, searchable plans for premium users**
  - This also allows user interaction, such as comments, likes and tagging.
- **Optional sign-up for a leaderboard of most active users**
- **Add usage statistics dashboard for each premium user**
  - KM swam in plans
  - Number of plans exported
  - Number of plans donated
  - Comparative statistics with other user behaviour
- **Add multi-plan generation as a one-time payment or premium feature**
  - Generate weeks or months of consecutive training plans as background task to prepare for a competition or to reach certain goals
- **Add premium model usage as a premium feature**
- **Add unlimited number of plans as a premium feature**
- **Add unlimited number of interactions as a premium feature**
  - Fee-tier users get 3 back-and-forth for a plan for free
- **TP are exportable as PDF, for premium users excel, csv or odt**
