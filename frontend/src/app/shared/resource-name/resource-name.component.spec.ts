import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ResourceNameComponent } from './resource-name.component';

describe('ResourceNameComponent', () => {
  let component: ResourceNameComponent;
  let fixture: ComponentFixture<ResourceNameComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ResourceNameComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceNameComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
