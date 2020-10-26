import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateOrEditGroupComponent } from './create-or-edit-group.component';

describe('CreateOrEditGroupComponent', () => {
  let component: CreateOrEditGroupComponent;
  let fixture: ComponentFixture<CreateOrEditGroupComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CreateOrEditGroupComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateOrEditGroupComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
