import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GroupOrUserPickerComponent } from './group-or-user-picker.component';

describe('GroupOrUserPickerComponent', () => {
  let component: GroupOrUserPickerComponent;
  let fixture: ComponentFixture<GroupOrUserPickerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ GroupOrUserPickerComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(GroupOrUserPickerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
