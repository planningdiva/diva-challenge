# diva-challenge Coding Assignment

This coding assignment is meant to test basic skills regarding the handling of auth tokens and API calls from a React or Next.js app.
It is a straightforward project that is meant to test basic proficiency with server communications that we work with daily at Planning Diva.
With that said, we are a new startup and creativity is important in our work.  Please feel free to add your own unique, creative ideas to this
project!

Time estimate:  1-4 hours.

Instructions and guidelines:
1. Create a React or Next.js project from scratch
2. Build a re-usable component and render it to the screen.  This component should include:
 - A button to check the status of the server.
 - A text field that contains a message.
 - A button to send a message to our engineering team, using the text from the above text field.
 - An area to display the response status from the last action.
3. Hook up the API routes provided below to make the above buttons and text field functional.
4. There should be no login button--A user should use the current JWT.  If the JWT is expired and the call fails, your application should retrieve a
new JWT behind the scenes and the user should never know the difference.
5. When complete, place your code up somewhere so the team can take a look at your solution.
6. Please reach out to josh@planningdiva.com for any questions or issues.


You can find the challenge server at this address:  `https://diva-challenge-ul4cm77qva-uc.a.run.app` with the following available endpoints:

GET `/login`:  A call to this endpoint will return a JSON Web Token (JWT).  The token has a 15-second expiration.  No payload needs to be sent.

GET `/alive`:  A call to this endpoint will return the text "alive" if the server is up.  This endpoint requires authentication with a JWT.  Include
a header in the request, with the following format:  "Authorization: Bearer JWT" where JWT is the token received from the `/login` endpoint.

POST `/slack`:  A call to this endpoint will post a message in our Planning Diva Slack space, so we can see that the coding assignment was successful.
This endpoint also requires authentication and should include a header with the same format as `/alive`.  This endpoint should include a body with the following format:  `{"text": "Place a message here"}`
 
 Good luck, and have fun!
