import {TestBed} from '@angular/core/testing';

import {BlocksService} from './blocks.service';
import {BackendService} from '../api/backend.service';
import {HttpClientTestingModule} from '@angular/common/http/testing';
import {HttpClientModule} from '@angular/common/http';
import {Message} from '../api/models';

describe('BlockService', () => {
  let service: BlocksService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule, HttpClientModule],
      providers: [BlocksService],
    });
    service = TestBed.inject(BlocksService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should have no message to start with', () => {
    expect(service.getMessage()).toBeUndefined();
  });

  it('should store the message', () => {
    service.setMessage({id: '1'} as Message);
    expect(service.getMessage().id).toBe('1');
  });

});
