import {Component, Input, OnInit} from '@angular/core';
import {ButtonElement} from '../../api/models';

@Component({
  selector: 'app-button',
  templateUrl: './button.component.html',
  styleUrls: ['./button.component.css']
})
export class ButtonComponent implements OnInit {

  constructor() {
  }

  @Input()
  buttonElement: ButtonElement;

  ngOnInit(): void {
  }

}
