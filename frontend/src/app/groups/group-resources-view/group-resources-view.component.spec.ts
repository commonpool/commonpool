import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GroupResourcesViewComponent } from './group-resources-view.component';

describe('GroupResourcesViewComponent', () => {
  let component: GroupResourcesViewComponent;
  let fixture: ComponentFixture<GroupResourcesViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ GroupResourcesViewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GroupResourcesViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
