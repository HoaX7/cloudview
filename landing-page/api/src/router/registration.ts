import { IRequest } from "itty-router";
import { Env } from "..";
import { Registration as RegistrationController } from "../controllers/registrations";
import { ThrowableRouter, withContent } from "itty-router-extras";

export function Register(router: ThrowableRouter) {
	router.post("/register", withContent, RegistrationController.create);
}
