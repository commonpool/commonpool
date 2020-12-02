it('loads examples', () => {
  cy.visit('/');
  cy.contains('Commonpool');
  cy.contains('In Construction.');
});
