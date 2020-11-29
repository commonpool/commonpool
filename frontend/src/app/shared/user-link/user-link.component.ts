import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-user-link',
  templateUrl: './user-link.component.html',
  styleUrls: ['./user-link.component.css']
})
export class UserLinkComponent implements OnInit {

  constructor() {
    console.log('new userlink');
  }

  @Input()
  id: string;

  ngOnInit(): void {
  }

}
