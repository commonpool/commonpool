import {Component, Input, OnInit} from '@angular/core';
import {Block} from '../../api/models';

@Component({
  selector: 'app-block',
  templateUrl: './block.component.html',
  styleUrls: ['./block.component.css']
})
export class BlockComponent implements OnInit {

  constructor() {
  }

  @Input()
  block: Block;

  ngOnInit(): void {
  }

}
