import {TestBed} from '@angular/core/testing';

import {AuthService} from './auth.service';
import {BackendService} from './api/backend.service';
import {Router} from '@angular/router';

describe('AuthService', () => {
  let service: AuthService;

  const mockBackend = {};
  const mockRouter = {};

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        {
          provide: BackendService,
          useValue: mockBackend
        }, {
          provide: Router,
          useValue: mockRouter
        }
      ]
    });
    service = TestBed.inject(AuthService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
