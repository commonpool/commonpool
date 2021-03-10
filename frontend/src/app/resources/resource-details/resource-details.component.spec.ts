import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {ResourceDetailsComponent} from './resource-details.component';
import {BackendService} from '../../api/backend.service';
import {ActivatedRoute, convertToParamMap, Router} from '@angular/router';
import {AuthService} from '../../auth.service';
import {of} from 'rxjs';
import {ExtendedResource, Resource, CallType} from '../../api/models';

describe('ResourceDetailsComponent', () => {
  let component: ResourceDetailsComponent;
  let fixture: ComponentFixture<ResourceDetailsComponent>;

  const paramMap = of(convertToParamMap({
    id: 'd31e8b48-7309-4c83-9884-4142efdf7271',
  }));

  const mockBackend = {
    getResource: () => {
      return of(
        new ExtendedResource(new Resource(
          'id',
          'summary',
          'description',
          CallType.Offer,
          10,
          20,
          'createdBy',
          'createdById',
          '2020-10-20T19:26:05.425Z',
          []
        )));
    }
  };

  const mockRouter = {};
  const mockRoute = {
    params: paramMap
  };
  const mockAuth = {};

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ResourceDetailsComponent],
      providers: [
        {
          provide: BackendService,
          useValue: mockBackend
        }, {
          provide: Router,
          useValue: mockRouter
        }, {
          provide: ActivatedRoute,
          useValue: mockRoute
        }, {
          provide: AuthService,
          useValue: mockAuth
        },
      ]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
})
;
