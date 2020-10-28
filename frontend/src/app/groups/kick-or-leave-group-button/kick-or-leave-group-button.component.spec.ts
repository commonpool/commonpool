import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KickOrLeaveGroupButtonComponent } from './kick-or-leave-group-button.component';

describe('KickOrLeaveGroupButtonComponent', () => {
  let component: KickOrLeaveGroupButtonComponent;
  let fixture: ComponentFixture<KickOrLeaveGroupButtonComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KickOrLeaveGroupButtonComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KickOrLeaveGroupButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
