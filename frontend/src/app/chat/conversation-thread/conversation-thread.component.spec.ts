import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ConversationThreadComponent } from './conversation-thread.component';

describe('ConversationThreadComponent', () => {
  let component: ConversationThreadComponent;
  let fixture: ComponentFixture<ConversationThreadComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ConversationThreadComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ConversationThreadComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
