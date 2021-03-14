import {ComponentFixture, fakeAsync, TestBed, tick} from '@angular/core/testing';
import {ValueComponent} from './value.component';
import {Component, ViewChild} from '@angular/core';
import {FormBuilder, FormsModule, ReactiveFormsModule} from '@angular/forms';
import {ValueRange} from '../api/models';
import {By} from '@angular/platform-browser';

@Component({
  template: `
    <div [formGroup]="form">
      <app-value formControlName="value"></app-value>
    </div>`
})
export class TestValueComponent {
  @ViewChild(ValueComponent)
  instance: ValueComponent;
  fb = new FormBuilder();
  form = this.fb.group({
    value: this.fb.control(undefined)
  });
}

describe('ValueComponent', () => {
  let component: TestValueComponent;
  let fixture: ComponentFixture<TestValueComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ReactiveFormsModule, FormsModule],
      declarations: [TestValueComponent, ValueComponent]
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TestValueComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should have null value by default', () => {
    expect(component.form.value.value).toBe(null);
  });

  it('should set the value', () => {
    component.form.setValue({value: new ValueRange(-1, 1)});
    expect(component.form.get('value').value.from).toBe(-1);
    expect(component.form.get('value').value.to).toBe(1);
  });

  it('should update value', () => {
    component.form.setValue({value: new ValueRange(-1, 1)});
    expect(component.form.get('value').value.from).toBe(-1);
    expect(component.form.get('value').value.to).toBe(1);
    component.form.setValue({value: new ValueRange(2, 3)});
    expect(component.form.get('value').value.from).toBe(2);
    expect(component.form.get('value').value.to).toBe(3);
  });

  it('should update value on input change', ((done) => {
    component.instance.min = -2;
    component.instance.max = 2;
    fixture.detectChanges();
    const element = fixture.debugElement.query(By.css('input')).nativeElement;
    element.value = '2';
    element.dispatchEvent(new Event('input'));
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(component.form.get('value').value.from).toBe(2);
      done();
    });
  }));

});
