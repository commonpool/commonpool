import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RequiredIndicatorComponent } from './required-indicator.component';

describe('RequiredIndicatorComponent', () => {
  let component: RequiredIndicatorComponent;
  let fixture: ComponentFixture<RequiredIndicatorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RequiredIndicatorComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RequiredIndicatorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
