// ***********************************************
// This example namespace declaration will help
// with Intellisense and code completion in your
// IDE or Text Editor.
// ***********************************************

import Chainable = Cypress.Chainable;
import RequestOptions = Cypress.RequestOptions;

// tslint:disable-next-line:no-namespace
declare namespace Cypress {
  interface Chainable<Subject = any> {
    login: typeof login;
  }
}
//
// function customCommand(param: any): void {
//   console.warn(param);
// }
//
// NOTE: You can use it like so:
// Cypress.Commands.add('customCommand', customCommand);
//
// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add("login", (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add("drag", { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add("dismiss", { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite("visit", (originalFn, url, options) => { ... })

Cypress.Commands.add('login', login);

function login(): Chainable<Cypress.AUTWindow> {
  Cypress.log({
    name: 'KeycloakLogin',
  });
  const options: Partial<RequestOptions> = {
    method: 'POST',
    url: Cypress.env('auth_url'),
    body: {
      grant_type: 'password',
      username: Cypress.env('auth_username'),
      password: Cypress.env('auth_password'),
      audience: Cypress.env('auth_audience'),
      scope: 'openid profile email',
      client_id: Cypress.env('auth_client_id'),
      client_secret: Cypress.env('auth_client_secret'),
    },
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded'
    }
  };
  return cy.request(options).then(({body}) => {
    const {access_token, refresh_token, expires_in, id_token} = body;
    const state = 'eyJkZXMiOiAiLyIsICJzdGF0ZSI6ICJzb21lU3RhdGUifQ==';
    const callbackUrl = `/api/v1/oauth2/callback?token=${access_token}&refresh_token=${refresh_token}&scope=openid&id_token=${id_token}&expires_in=${expires_in}&token_type=Bearer&state=${state}`;
    return cy.visit(callbackUrl, {});
  });
}
