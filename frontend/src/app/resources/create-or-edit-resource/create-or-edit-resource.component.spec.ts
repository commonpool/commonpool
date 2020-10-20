import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateOrEditResourceComponent } from './create-or-edit-resource.component';

describe('NewResourceComponent', () => {
  let component: CreateOrEditResourceComponent;
  let fixture: ComponentFixture<CreateOrEditResourceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CreateOrEditResourceComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateOrEditResourceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
