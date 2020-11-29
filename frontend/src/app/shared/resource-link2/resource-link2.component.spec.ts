import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ResourceLink2Component } from './resource-link2.component';

describe('ResourceLink2Component', () => {
  let component: ResourceLink2Component;
  let fixture: ComponentFixture<ResourceLink2Component>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ResourceLink2Component ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceLink2Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
