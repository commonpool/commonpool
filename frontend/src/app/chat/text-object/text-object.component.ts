import {Component, Input, OnInit, ViewEncapsulation} from '@angular/core';
import {TextObject} from '../../api/models';

@Component({
  selector: 'app-text-object',
  templateUrl: './text-object.component.html',
  styleUrls: ['./text-object.component.css'],
  encapsulation: ViewEncapsulation.None
})
export class TextObjectComponent implements OnInit {

  @Input()
  textObject: TextObject;

  @Input()
  subtle = false;

  @Input()
  small = false;

  constructor() {
  }

  ngOnInit(): void {
  }

}
