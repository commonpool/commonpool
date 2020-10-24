import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ResourceInquiryComponent } from './resource-inquiry.component';

describe('ResourceInquiryComponent', () => {
  let component: ResourceInquiryComponent;
  let fixture: ComponentFixture<ResourceInquiryComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ResourceInquiryComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceInquiryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
