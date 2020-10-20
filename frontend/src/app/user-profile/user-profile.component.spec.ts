import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {UserProfileComponent} from './user-profile.component';
import {ActivatedRoute, convertToParamMap} from '@angular/router';
import {BackendService} from '../api/backend.service';
import {of} from 'rxjs';
import {UserInfoResponse} from '../api/models';

describe('UserProfileComponent', () => {
  let component: UserProfileComponent;
  let fixture: ComponentFixture<UserProfileComponent>;

  const paramMap = of(convertToParamMap({
    id: 'd31e8b48-7309-4c83-9884-4142efdf7271',
  }));

  const mockActivatedRoute = {
    params: paramMap
  };

  const mockBackend = {
    getUserInfo: () => of({
      id: 'user-id',
      username: 'username'
    } as UserInfoResponse),
    searchResources: () => of([])
  };

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [UserProfileComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: mockActivatedRoute
        }, {
          provide: BackendService,
          useValue: mockBackend
        }
      ]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserProfileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
