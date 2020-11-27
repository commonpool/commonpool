import {Component, Input} from '@angular/core';
import {Block} from '../../api/models';

@Component({
  selector: 'app-context-block',
  template: `
    <div class="d-flex block block-context">
      <ng-container *ngFor="let e of block.elements; let i = index">
        <div [class.pl-1]="i !== 0">
          <ng-container *ngIf="e.type === 'image'">
            <img style="max-height: 1.25rem" class="d-inline-block" [src]="e.imageUrl"
                 [alt]="e.altText">
          </ng-container>
          <ng-container *ngIf="e.type === 'text' || e.type === 'mrkdwn'">
            <app-text-object [textObject]="e" [small]="true"></app-text-object>
          </ng-container>
        </div>
      </ng-container>
    </div>`,
  styles: [`
    .block-context {
      font-size: 0.75rem;
      color: gray;
    }
  `]
})
export class ContextBlockComponent {

  constructor() {
  }

  @Input()
  block: Block;

}
