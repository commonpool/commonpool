import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GroupMemberViewComponent } from './group-member-view.component';

describe('GroupMemberViewComponent', () => {
  let component: GroupMemberViewComponent;
  let fixture: ComponentFixture<GroupMemberViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ GroupMemberViewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GroupMemberViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
