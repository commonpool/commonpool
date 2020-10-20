import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {CreateOrEditResourceComponent} from './create-or-edit-resource.component';
import {BackendService} from '../../api/backend.service';
import {ActivatedRoute, convertToParamMap} from '@angular/router';
import {AuthService} from '../../auth.service';
import {of} from 'rxjs';

describe('CreateOrEditResourceComponent', () => {
  let component: CreateOrEditResourceComponent;
  let fixture: ComponentFixture<CreateOrEditResourceComponent>;

  const paramMap = of(convertToParamMap({
    id: 'd31e8b48-7309-4c83-9884-4142efdf7271',
  }));

  const mockBackend = {};
  const mockRoute = {
    params: paramMap
  };
  const mockAuth = {};

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [CreateOrEditResourceComponent],
      providers: [
        {
          provide: BackendService,
          useValue: mockBackend
        }, {
          provide: ActivatedRoute,
          useValue: mockRoute,
        }, {
          provide: AuthService,
          useValue: mockAuth
        }
      ]
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
