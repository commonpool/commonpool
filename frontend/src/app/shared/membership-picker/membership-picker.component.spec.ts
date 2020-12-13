import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MembershipPickerComponent } from './membership-picker.component';

describe('MembershipPickerComponent', () => {
  let component: MembershipPickerComponent;
  let fixture: ComponentFixture<MembershipPickerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MembershipPickerComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MembershipPickerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
