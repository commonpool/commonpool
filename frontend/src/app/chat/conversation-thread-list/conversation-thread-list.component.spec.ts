import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {ConversationThreadListComponent} from './conversation-thread-list.component';
import {of} from 'rxjs';
import {BackendService} from '../../api/backend.service';

describe('ConversationThreadListComponent', () => {
  let component: ConversationThreadListComponent;
  let fixture: ComponentFixture<ConversationThreadListComponent>;

  const mockBackend = {
    getThreads: () => of([])
  };
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ConversationThreadListComponent],
      providers: [
        {
          provide: BackendService,
          useValue: mockBackend
        }
      ]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ConversationThreadListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
