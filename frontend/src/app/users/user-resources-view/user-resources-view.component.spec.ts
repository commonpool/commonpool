import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserResourcesViewComponent } from './user-resources-view.component';

describe('UserResourcesViewComponent', () => {
  let component: UserResourcesViewComponent;
  let fixture: ComponentFixture<UserResourcesViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UserResourcesViewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserResourcesViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
