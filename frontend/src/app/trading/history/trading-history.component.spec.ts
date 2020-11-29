import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TradingHistoryComponent } from './trading-history.component';

describe('TradingHistoryComponent', () => {
  let component: TradingHistoryComponent;
  let fixture: ComponentFixture<TradingHistoryComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TradingHistoryComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TradingHistoryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
