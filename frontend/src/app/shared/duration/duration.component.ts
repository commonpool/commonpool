import {Component, Input} from '@angular/core';

@Component({
  selector: 'app-duration',
  template: `{{duration / 3600000000000}}h`,
})
export class DurationComponent {
  @Input()
  duration: number;
}
