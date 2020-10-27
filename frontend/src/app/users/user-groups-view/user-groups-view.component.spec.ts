import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserGroupsViewComponent } from './user-groups-view.component';

describe('UserGroupsViewComponent', () => {
  let component: UserGroupsViewComponent;
  let fixture: ComponentFixture<UserGroupsViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UserGroupsViewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserGroupsViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
