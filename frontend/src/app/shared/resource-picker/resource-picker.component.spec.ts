import {async, ComponentFixture, fakeAsync, flush, TestBed, waitForAsync} from '@angular/core/testing';

import {ResourcePickerComponent} from './resource-picker.component';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {NgSelectModule} from '@ng-select/ng-select';
import {BackendService} from '../../api/backend.service';
import {Component, ViewChild} from '@angular/core';
import {GetResourceResponse, Resource, CallType} from '../../api/models';
import {of} from 'rxjs';

@Component({
  template: `
    <app-resource-picker [(ngModel)]="resource"></app-resource-picker>
  `
})
export class TestResourcePickerComponent {

  @ViewChild(ResourcePickerComponent, {static: true})
  resourcePicker: ResourcePickerComponent;

  public resource: string | Resource;
}

describe('ResourcePickerComponent', () => {
  let component: TestResourcePickerComponent;
  let fixture: ComponentFixture<TestResourcePickerComponent>;

  const user1id = '6ec26bc4-b811-4e17-afae-c000d69956c5';

  const resource1 = new Resource(
    '7c851fde-ea2d-49a1-a867-0367fe3e8634',
    'Summary1',
    'Description1',
    CallType.Offer,
    1,
    2,
    'user1',
    user1id,
    new Date().toISOString(),
    []);

  const resource2 = new Resource(
    'c5f08bd1-338b-4a33-825f-359669bd510d',
    'Summary2',
    'Description2',
    CallType.Offer,
    2,
    2,
    'user1',
    user1id,
    new Date().toISOString(),
    []);

  const backend = {
    searchResources: jest.fn(),
    getResource: jest.fn()
  };

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TestResourcePickerComponent, ResourcePickerComponent],
      imports: [
        FormsModule,
        ReactiveFormsModule,
        NgSelectModule
      ],
      providers: [
        {provide: BackendService, useValue: backend}
      ]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TestResourcePickerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should have undefined value', waitForAsync(() => {
    const writeValue = spyOn(component.resourcePicker, 'writeValue').and.callThrough();
    expect(component.resource).toBeUndefined();
  }));

  it('should query the backend for the resource when provided with an id', fakeAsync(() => {
    backend.getResource.mockReturnValueOnce(of(new GetResourceResponse(resource1)));
    component.resource = resource1.resourceId;
    fixture.detectChanges();
    flush();
    expect(backend.getResource).toHaveBeenCalledWith(resource1.resourceId);
  }));

  it('should not query the backend for the resource when provided with a resource', fakeAsync(() => {
    backend.getResource.mockReturnValueOnce(of(new GetResourceResponse(resource1)));
    component.resource = resource1;
    fixture.detectChanges();
    flush();
    expect(backend.getResource).not.toHaveBeenCalled();
  }));

});
