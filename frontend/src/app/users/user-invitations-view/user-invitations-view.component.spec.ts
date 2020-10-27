import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserInvitationsViewComponent } from './user-invitations-view.component';

describe('UserInvitationsViewComponent', () => {
  let component: UserInvitationsViewComponent;
  let fixture: ComponentFixture<UserInvitationsViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UserInvitationsViewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserInvitationsViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
