package garmin

// TODO: implement the garmin api, based on the unofficial api, but only for workouts as this is the scope of this project atm
// Use python implementation as reference, focused on swim workouts
// https://github.com/matin/garth/blob/main/src/garth/http.py
// https://github.com/cyberjunky/python-garminconnect/blob/master/garminconnect/workout.py
// def upload_workout(
//     self, workout_json: dict[str, Any] | list[Any] | str
// ) -> dict[str, Any]:
//     """Upload workout using json data."""
//     url = f"{self.garmin_workouts}/workout"
//     logger.debug("Uploading workout using %s", url)

//     if isinstance(workout_json, str):
//         import json as _json

//         try:
//             payload = _json.loads(workout_json)
//         except Exception as e:
//             raise ValueError(f"invalid workout_json string: {e}") from e
//     else:
//         payload = workout_json
//     if not isinstance(payload, dict | list):
//         raise ValueError("workout_json must be a JSON object or array")
//     return self.garth.post("connectapi", url, json=payload, api=True).json()
