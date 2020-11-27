import {ComponentFixture, TestBed} from '@angular/core/testing';

import {ContextBlockComponent} from './context-block.component';

describe('ContextBlockComponent', () => {
  let component: ContextBlockComponent;
  let fixture: ComponentFixture<ContextBlockComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ContextBlockComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ContextBlockComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
