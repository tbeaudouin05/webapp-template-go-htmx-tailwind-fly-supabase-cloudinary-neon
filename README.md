should search and replace: webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon
by your project name everywhere in the code, in particular go.mod AND fly.toml.

Should install Go: go version go1.23.0 darwin/arm64

Run: go get to install packages in the project

Install nodejs (for Tailwind)



Should create new:
- fly.io project
- Neon project
- supabase project
- cloudinary project


Should connect repository to your fly.io project
Should update env variables for neon, supabase, cloudinary access

fly.io:
1. Install: https://fly.io/docs/flyctl/install/
2. Run: fly launch
3. IMPORTANT: Would you like to copy its configuration to the new app? --> SAY YES!


Should create file .env at root, find structure of .env in struct EnvVar in file goEnv.go
Please note that environment variable names should have exact same name in EnvVar as in .env
Example: 
- Env should be called Env in .env file, not ENV
- ShouldUseCdn should be called as such and not SHOULD_USE_CDN

Generally speaking, we only use camel case everywhere.

To add tailwind animations, use tailwind.config.js (ask chatgpt how to do this with tailwind.config.js) - animations will only be available if ShouldUseCdn = No

Static files can be served from frontend/static/

To run locally, go to cmd/run and run: go run .

To deploy, go to cmd/deploy and run: go run .


