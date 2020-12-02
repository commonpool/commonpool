describe('create resource', () => {

  beforeEach(() => {
    cy.login();
    cy.visit('/resources/new');
  });

  it('should not allow negative time values', () => {



    cy.get('[data-cy=valueInHoursTo]')
      .debug()
      .clear()
      .type('-1')
      .blur()
      .should('have.value', 0);

    cy.get('[data-cy=valueInHoursFrom]')
      .debug()
      .clear()
      .type('-3')
      .blur();


    cy.get('[data-cy=valueInHoursFrom]')
      .should('have.value', 0);


  });

  it('should allow to create a resource', () => {
    cy
      .get('[data-cy=summary]')
      .type('A Summary')
      .get('[data-cy=description]')
      .type('A Description')
      .get('[data-cy=valueInHoursFrom]')
      .clear()
      .type('2')
      .get('[data-cy=valueInHoursTo]')
      .clear()
      .type('10')
      .get('[data-cy=submit]')
      .click();

    cy.get('app-resource-details')
      .contains('A Summary');

    cy.get('app-resource-details')
      .contains('A Description');

  });

});
