import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserOffersViewComponent } from './user-offers-view.component';

describe('UserOffersViewComponent', () => {
  let component: UserOffersViewComponent;
  let fixture: ComponentFixture<UserOffersViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UserOffersViewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserOffersViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
