import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-resource-link2',
  templateUrl: './resource-link2.component.html',
  styleUrls: ['./resource-link2.component.css']
})
export class ResourceLink2Component implements OnInit {

  constructor() {
  }

  @Input()
  id: string;

  ngOnInit(): void {
  }

}
