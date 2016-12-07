/* @flow */
import {AuthenticationForm, SignIn, SignUp, ForgotPassword} from './authForms';

function main(): void{
    let signInForm: HTMLFormElement = ((document.getElementById('auth-form'):any):HTMLFormElement);

    /*
     * Set up authentication forms
     */
    if(signInForm != null){
        try{
            let formHandler: AuthenticationForm = authFormHandler(signInForm);
            signInForm.addEventListener('submit', function(event: Event): void{
                event.preventDefault();
                formHandler.submit();
            });
        }catch(error){
            console.log(error.message);
        }
    }
}

function authFormHandler(form: HTMLFormElement): AuthenticationForm{
    let formName: string = form.getAttribute('name');
    switch(formName){
        case 'signin':
            return new SignIn(form);
        case 'signup':
            return new SignUp(form);
        case 'forgot':
            return new ForgotPassword(form);
        default:
            throw new Error('Missing valid name attribute on auth-form');
    }
    
}

// Run main
document.addEventListener('DOMContentLoaded', main);
