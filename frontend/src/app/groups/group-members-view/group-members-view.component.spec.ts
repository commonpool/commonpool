import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GroupMembersViewComponent } from './group-members-view.component';

describe('GroupMembersViewComponent', () => {
  let component: GroupMembersViewComponent;
  let fixture: ComponentFixture<GroupMembersViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ GroupMembersViewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GroupMembersViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
