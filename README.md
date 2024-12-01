should search and replace: webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon
by your project name everywhere in the code, in particular go.mod AND fly.toml.

Should install Go: go version go1.23.0 darwin/arm64

Run: go get to install packages in the project

Install nodejs (for Tailwind)



Should create new:
- fly.io project
- passage account
- neon project
- cloudinary project --> to do
- example complex form with https://github.com/Emtyloc/json-enc-custom/blob/main/README.md
- example upload image cloudinary
- stripe?


Should connect repository to your fly.io project:
1. Install: https://fly.io/docs/flyctl/install/
2. Run: fly launch
3. IMPORTANT: Would you like to copy its configuration to the new app? --> SAY YES!


Should update env variables for neon, supabase, cloudinary access


Should create file .env at root, find structure of .env in struct EnvVar in file goEnv.go
Please note that environment variable names should have exact same name in EnvVar as in .env
Example: 
- Env should be called Env in .env file, not ENV
- ShouldUseCdn should be called as such and not SHOULD_USE_CDN

Generally speaking, we only use camel case everywhere.



----
Create neon project here: https://console.neon.tech/app/projects

get your database url here: https://console.neon.tech/app/projects/falling-fog-13122533/quickstart

update NeonDatabaseUrl in .env file with your database url

(rename env variable DATABASE_URL to NeonDatabaseUrl in .env file)

----

Create Passage account here: https://console.passage.id/register

create a new Passage app here: https://console.passage.id/
IMPORTANT: for local testing, you should set the authentication origin to: http://0.0.0.0:8080
For non-local environment, set the authentication origin to the host address (ex: https://www.example.com)
This means that you should create different Passage apps for local vs. non-local environments.

find your PassageAppId here once you created your passage app: https://console.passage.id/

update PassageAppId in .env file

use this to logout user:
import (
  "net/http"
 
  "github.com/passageidentity/passage-go"
)
 
func exampleHandler(w http.ResponseWriter, r *http.Request) {
 
  psg, _ := passage.New("<PASSAGE_APP_ID>", &passage.Config{
    APIKey: "<PASSAGE_API_KEY>",
  })
 
  // Get the Passage User ID from database
 
  // The passageUserID returned above can be used deactivate a user
  deactivatedUser, err := psg.DeactivateUser(passageUserID)
  if err != nil {
    // Couldn't deactivate the user
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
 
  // The passageUserID returned above can be used activate a user
  activatedUser, err := psg.ActivateUser(passageUserID)
  if err != nil {
    // Couldn't activate the user
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
 
}

----

To add tailwind animations, use tailwind.config.js (ask chatgpt how to do this with tailwind.config.js) - animations will only be available if ShouldUseCdn = No

Static files can be served from frontend/static/

To run locally, go to cmd/run and run: go run .

To deploy, go to cmd/deploy and run: go run .


