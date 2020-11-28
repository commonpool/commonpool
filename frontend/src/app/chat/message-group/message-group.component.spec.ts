import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageGroupComponent } from './message-group.component';

describe('MessageGroupComponent', () => {
  let component: MessageGroupComponent;
  let fixture: ComponentFixture<MessageGroupComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageGroupComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageGroupComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
