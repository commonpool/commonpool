describe('login', () => {
  it('should successfully log into our app', () => {
    cy.login().then(() => {
      cy.visit('/');
    });
  });
});

