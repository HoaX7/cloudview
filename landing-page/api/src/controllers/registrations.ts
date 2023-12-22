import { jsonError, jsonSuccess } from "../helper/response";
import { IRequest } from "itty-router";
import { Env } from "..";
import { ALREADY_EXISTS_ERROR, INTERNAL_SERVER_ERROR, INVALID_EMAIL_ERROR, INVALID_USERNAME_ERROR } from "../helper/errors";
import { validateEmail } from "../helper/email";

export const Registration = {
	async create({ content }: IRequest, env: Env) {
		console.log(`Content: `, content);
		try {
			const email: string = content.email?.trim();
			const name: string = content.name?.trim();
			if (!validateEmail(email)) {
				return jsonError(
					{
						success: false,
						data: { message: INVALID_EMAIL_ERROR.error, code: INVALID_EMAIL_ERROR.code },
					},
					{ status: 422 }
				);
			}
			if (!name) {
				return jsonError(
					{
						success: false,
						data: { message: INVALID_USERNAME_ERROR.error, code: INVALID_USERNAME_ERROR.code }
					},
					{ status: 422 }
				)
			}

			const { results } = await env.HVEC_MARKETING_DB.prepare(
				"INSERT INTO registrations (email_address, name) VALUES (?, ?) RETURNING *"
			)
				.bind(email, name)
				.all();
			console.log(`Inserted result`, JSON.stringify(results));
			return jsonSuccess({});
		} catch (err) {
			const error = err as Error & { cause: { message: string } };
			console.error("Error while inserting registrations email", { msg: error.cause?.message || error.message });
			console.error("error.message", error.message)
			if (error.message?.includes("D1_ERROR") && error.cause?.message?.toLowerCase().includes("unique constraint failed")) {
				return jsonError({ message: ALREADY_EXISTS_ERROR.error, code: ALREADY_EXISTS_ERROR.code }, { status: 422 });
			}
			return jsonError({ message: INTERNAL_SERVER_ERROR.error, code: INTERNAL_SERVER_ERROR.code }, { status: 500 });
		}
	},
};
