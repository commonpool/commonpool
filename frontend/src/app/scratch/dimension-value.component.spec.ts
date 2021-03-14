import {ComponentFixture, TestBed} from '@angular/core/testing';
import {Component, ViewChild} from '@angular/core';
import {FormBuilder, FormsModule, ReactiveFormsModule} from '@angular/forms';
import {DimensionValueComponent} from './dimension-value.component';
import {ValueComponent} from './value.component';
import {ValueThresholdComponent} from './value-threshold.component';
import {By} from '@angular/platform-browser';
import {ValueDimension, ValueRange} from '../api/models';
import {ValueDimensionService} from './dimension.service';
import {of} from 'rxjs';

@Component({
  template: `
    <div [formGroup]="form">
      <app-dimension-value formControlName="value"></app-dimension-value>
    </div>`
})
export class TestDimensionValueComponent {
  @ViewChild(DimensionValueComponent)
  instance: DimensionValueComponent;
  fb = new FormBuilder();
  form = this.fb.group({
    value: this.fb.control(undefined)
  });
}

describe('DimensionValueComponent', () => {
  let component: TestDimensionValueComponent;
  let fixture: ComponentFixture<TestDimensionValueComponent>;

  const dimension = new ValueDimension(
    'test',
    'summary',
    1,
    new ValueRange(-3, 4),
    [{description: 'threshold1'}, {description: 'threshold2'}]);

  const mockSvc = {
    getDimension: () => {
      return of(dimension);
    }
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ReactiveFormsModule, FormsModule],
      providers: [
        {
          provide: ValueDimensionService,
          useValue: mockSvc
        }
      ],
      declarations: [
        ValueThresholdComponent,
        ValueComponent,
        TestDimensionValueComponent,
        DimensionValueComponent]
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TestDimensionValueComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should have null value', () => {
    expect(component.form.value.value).toBeNull();
  });

  it('should update value on input change', ((done) => {

    component.form.setValue({
      value: {
        dimensionName: 'test',
        valueRange: {
          from: -1,
          to: 2
        }
      }
    });

    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(component.form.get('value').value.valueRange.from).toBe(-1);
      expect(component.form.get('value').value.valueRange.to).toBe(2);
      expect(component.form.get('value').value.dimensionName).toBe('test');
    }).then(() => {
      const element = fixture.debugElement.query(By.css('input')).nativeElement;
      element.value = '3';
      element.dispatchEvent(new Event('input'));
      fixture.detectChanges();
      return fixture.whenStable();
    }).then(() => {
      expect(component.form.get('value').value.valueRange.from).toBe(3);
      expect(component.form.get('value').value.valueRange.to).toBe(3);
      expect(component.form.get('value').value.dimensionName).toBe('test');
      done();
    });
  }));

});
