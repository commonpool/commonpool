import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GroupInvitesViewComponent } from './group-invites-view.component';

describe('GroupInvitesViewComponent', () => {
  let component: GroupInvitesViewComponent;
  let fixture: ComponentFixture<GroupInvitesViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ GroupInvitesViewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GroupInvitesViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
