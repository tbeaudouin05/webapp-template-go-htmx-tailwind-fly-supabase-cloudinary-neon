should search and replace: webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon
by your project name everywhere in the code, in particular go.mod AND fly.toml.

Should install Go: go version go1.23.0 darwin/arm64

Run: go get to install packages in the project

Install nodejs & npm (for Tailwind)

---

For creating web components in Go, we use:
- Templ: https://templ.guide/ (this is a Go package for creating HTML components in Go)
- HTMX: https://htmx.org/docs/ (this is a lightweight javascript library that extends HTML in order to, among other things, enable to replace only part of the HTML of the page when the server sends back HTML)


---

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

----

To add tailwind animations, use tailwind.config.js (ask chatgpt how to do this with tailwind.config.js) - animations will only be available if ShouldUseCdn = No

By default, this project is set up to use Tailwind Flowbite which is a tailwind library for pre-built design components - but you can also use plain tailwind or any other tailwind library. Flowbite is the Tailwind plugin with the most comprehensive list of pre-built components (more than 400 components).
If you do not want to use Flowbite components, just delete this: require('flowbite/plugin'), in tailwind.config.js

If you want to use Flowbite, you should run: npm install flowbite to install it.

--

Static files can be served from frontend/static/

To run locally, go to cmd/run and run: go run .

To deploy, go to cmd/deploy and run: go run .
Note: before deploying, you should set the fly env variables.
For each variable in .env, you should run: fly secrets set VariableName=VariableValue
For instance: 
- fly secrets set Env=development
- fly secrets set ShouldUseCdn=No
Same for NeonDatabaseUrl and PassageAppId


-- 

For help to code with flowbite, the flowbite team created a flowbite GPT!
https://chatgpt.com/g/g-y7yC35HB9-flowbite-gpt


--

for "ORM" like feature: sqlc golang