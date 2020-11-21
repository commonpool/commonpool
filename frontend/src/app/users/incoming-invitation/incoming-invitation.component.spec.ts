import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IncomingInvitationComponent } from './incoming-invitation.component';

describe('IncomingInvitationComponent', () => {
  let component: IncomingInvitationComponent;
  let fixture: ComponentFixture<IncomingInvitationComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ IncomingInvitationComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IncomingInvitationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
