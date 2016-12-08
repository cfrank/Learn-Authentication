## Flow of the authentication

### Will be updated as code progresses

1. New authentication request (Signup)
   * A POST request is made to the path `/auth/signup`
  
     The request is made with the following data
     ```json
     {
    	 authString: <Base64 Encoded username:password string in that format>
         date: <Unix time in seconds>
         nonce: <Random string of bytes>
     }
     ```
   * The data is read into a struct of type `AuthenticationData`
   * The data is validated, including checking the date provided is not older than `MAXTIMEDIFF`
   * An authentication selector which will allow lookup in the `auth_token` database is created
   * The password is hashed with Scrypt along with a securly random salt
   * The auth_token is returned in the format `selector:hash`