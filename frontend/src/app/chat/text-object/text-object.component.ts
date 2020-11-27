import {Component, Input, ViewEncapsulation} from '@angular/core';
import {TextObject} from '../../api/models';

@Component({
  selector: 'app-text-object',
  template: `
    <span [class.text-muted]="subtle" [ngStyle]="{'font-size': small ? '0.875rem' : ''}">
        <ng-container *ngIf="textObject.type === 'plain_text'" style="white-space: pre-wrap">
            <span [innerText]="textObject.value"></span>
        </ng-container>
        <ng-container *ngIf="textObject.type === 'mrkdwn'">
            <markdown emoji [data]="textObject.value"></markdown>
        </ng-container>
    </span>
  `,
  styles: [`
    markdown :last-child {
      margin-bottom: 0;
    }
  `],
  encapsulation: ViewEncapsulation.None
})
export class TextObjectComponent {

  @Input()
  textObject: TextObject;

  @Input()
  subtle = false;

  @Input()
  small = false;

}
