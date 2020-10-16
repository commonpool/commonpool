import { TestBed } from '@angular/core/testing';

import { BackendService } from './backend.service';
import {HttpClientTestingModule} from '@angular/common/http/testing';

describe('BackendService', () => {
  let service: BackendService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule]
    });
    service = TestBed.inject(BackendService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
