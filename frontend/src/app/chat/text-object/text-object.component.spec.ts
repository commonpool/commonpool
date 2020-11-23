import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TextObjectComponent } from './text-object.component';

describe('TextObjectComponent', () => {
  let component: TextObjectComponent;
  let fixture: ComponentFixture<TextObjectComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TextObjectComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TextObjectComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
