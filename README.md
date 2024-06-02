# Go Monzo OAuth Flow

A basic third party [Monzo](https://monzo.com/) OAuth client for accessing the [Monzo API](https://docs.monzo.com/) written in Go.

**Note**

The access token and refresh token are currently not stored after the program executes. To persist this, you will need to store these tokens outside of the program. This has not been implemented for this project, so you will need to go through the authentication flow each time this program is run.

## Prerequisites

Before running the program, you'll need to register a new API client on the [Monzo Developer Portal](https://developers.monzo.com). 

To register a new API client, log inot the Monzo Devloper Portal (remember to approve the login via your Monzo app) and click "New OAuth Client".

Then provide the following details for your OAuth client (Logo URL can remain blank):

```
Name: Go Monzo Client  
Redirect URL: http://127.0.0.1:21234/callback
Description: Go Monzo Client Application
Confidentiality: True
```

Once the client is register you will recieve a Client ID and a Client Secert. Make a note of these!

## How to use
Clone the repository (I like using the GitHub CLI for this)
```bash
gh repo clone Skyth3r/go-monzo-client
```

Install dependencies
```bash
go mod tidy
```

Set Client ID and Client Secert in environment variables
```bash
export MONZO_CLIENT_ID=YOUR_CLIENT_ID_HERE

export MONZO_CLIENT_SECRET=YOUR_CLIENT_SECRET_HERE
```

Run the program
```bash
go run ./
```

## Expected results

This client will start the OAuth flow, and attempt to open a browser with the login URL. On the login page, type in your email address linked to your personal Monzo account and then click the link sent to your email address and go back to the app.

You will then be prompted to open the Monzo app and grant access to the app by clicking "Allow access to your data". This process is related to Strong Customer Authentication. Once access has been granted via the Monzo app, go back to the app and press the [Enter] key to continue.

The app will attempt to make an api request to the accounts endpoint and print out all the accounts you have with Monzo.

## Extending the client

This is a basic client as currently it does not store the retrieved access token and refresh token or refresh the access token and only makes one api call after completing Strong Customer Authentication.

It should be used as a foundation that can be build upon and expanded further.